#version 440 core

uniform mat4 modelViewProjectionMatrix;
uniform mat4 modelViewMatrix;

// ボーン変形行列を格納するテクスチャ
uniform sampler2D boneMatrixTexture;
uniform int boneMatrixWidth;
uniform int boneMatrixHeight;

in layout(location = 0) vec3 position;
in layout(location = 1) vec3 normal;
in layout(location = 2) vec2 uv;
in layout(location = 3) vec2 extendUv;
in layout(location = 4) float vertexEdge;
in layout(location = 5) vec4 boneIndexes;
in layout(location = 6) vec4 boneWeights;
in layout(location = 7) float isSdef;
in layout(location = 8) vec3 sdefC;
in layout(location = 9) vec3 sdefR0;
in layout(location = 10) vec3 sdefR1;
in layout(location = 11) vec3 sdefB0; // SDEF用ボーン0の位置
in layout(location = 12) vec3 sdefB1; // SDEF用ボーン1の位置

uniform vec4 diffuse;
uniform vec3 ambient;
uniform vec4 specular;

uniform vec3 cameraPosition;
uniform vec3 lightDirection;

uniform int useToon;
uniform int useSphere;
uniform int sphereMode;

out vec4 vertexColor;
out vec3 vertexSpecular;
out vec2 vertexUv;
out vec3 vetexNormal;
out vec2 sphereUv;
out vec3 eye;
out float totalBoneWeight;

// 球形補間
vec4 slerp(vec4 q1, vec4 q2, float t) {
    float dot = dot(q1, q2);

    if(dot < 0.0) {
        q1 = -q1; // q1の向きを反転させる
        dot = -dot;
    }

    if(dot > 0.9995) {
        // クォータニオンが非常に近い場合は線形補間を使用し、正規化する
        vec4 result = q1 + t * (q2 - q1);
        return normalize(result);
    }

    dot = clamp(dot, -1.0, 1.0); // 数値誤差による範囲外の値を修正
    float theta_0 = acos(dot); // q1とq2の間の角度
    float theta = theta_0 * t; // 現在のtにおける角度

    vec4 q3 = q2 - q1 * dot;
    q3 = normalize(q3); // 正規直交基底を作成

    return q1 * cos(theta) + q3 * sin(theta);
}

// 行列の逆行列を求める
mat4 inverseMatrix(mat4 m) {
    float s0 = m[0][0] * m[1][1] - m[1][0] * m[0][1];
    float s1 = m[0][0] * m[1][2] - m[1][0] * m[0][2];
    float s2 = m[0][0] * m[1][3] - m[1][0] * m[0][3];
    float s3 = m[0][1] * m[1][2] - m[1][1] * m[0][2];
    float s4 = m[0][1] * m[1][3] - m[1][1] * m[0][3];
    float s5 = m[0][2] * m[1][3] - m[1][2] * m[0][3];

    float c5 = m[2][2] * m[3][3] - m[3][2] * m[2][3];
    float c4 = m[2][1] * m[3][3] - m[3][1] * m[2][3];
    float c3 = m[2][1] * m[3][2] - m[3][1] * m[2][2];
    float c2 = m[2][0] * m[3][3] - m[3][0] * m[2][3];
    float c1 = m[2][0] * m[3][2] - m[3][0] * m[2][2];
    float c0 = m[2][0] * m[3][1] - m[3][0] * m[2][1];

    // Should check for 0 determinant
    float invdet = 1.0 / (s0 * c5 - s1 * c4 + s2 * c3 + s3 * c2 - s4 * c1 + s5 * c0);

    mat4 invM;

    invM[0][0] = (m[1][1] * c5 - m[1][2] * c4 + m[1][3] * c3) * invdet;
    invM[0][1] = (-m[0][1] * c5 + m[0][2] * c4 - m[0][3] * c3) * invdet;
    invM[0][2] = (m[3][1] * s5 - m[3][2] * s4 + m[3][3] * s3) * invdet;
    invM[0][3] = (-m[2][1] * s5 + m[2][2] * s4 - m[2][3] * s3) * invdet;

    invM[1][0] = (-m[1][0] * c5 + m[1][2] * c2 - m[1][3] * c1) * invdet;
    invM[1][1] = (m[0][0] * c5 - m[0][2] * c2 + m[0][3] * c1) * invdet;
    invM[1][2] = (-m[3][0] * s5 + m[3][2] * s2 - m[3][3] * s1) * invdet;
    invM[1][3] = (m[2][0] * s5 - m[2][2] * s2 + m[2][3] * s1) * invdet;

    invM[2][0] = (m[1][0] * c4 - m[1][1] * c2 + m[1][3] * c0) * invdet;
    invM[2][1] = (-m[0][0] * c4 + m[0][1] * c2 - m[0][3] * c0) * invdet;
    invM[2][2] = (m[3][0] * s4 - m[3][1] * s2 + m[3][3] * s0) * invdet;
    invM[2][3] = (-m[2][0] * s4 + m[2][1] * s2 - m[2][3] * s0) * invdet;

    invM[3][0] = (-m[1][0] * c3 + m[1][1] * c1 - m[1][2] * c0) * invdet;
    invM[3][1] = (m[0][0] * c3 - m[0][1] * c1 + m[0][2] * c0) * invdet;
    invM[3][2] = (-m[3][0] * s3 + m[3][1] * s1 - m[3][2] * s0) * invdet;
    invM[3][3] = (m[2][0] * s3 - m[2][1] * s1 + m[2][2] * s0) * invdet;

    return invM;
}

// クォータニオンvec4の逆回転を求める
vec4 inverseQuaternion(vec4 q) {
    return vec4(-q.x, -q.y, -q.z, q.w);
}

// 行列の転置行列を求める
mat4 transposeMatrix(mat4 m) {
    mat4 transposedMatrix;
    for(int i = 0; i < 4; ++i) {
        for(int j = 0; j < 4; ++j) {
            transposedMatrix[i][j] = m[j][i];
        }
    }
    return transposedMatrix;
}

// mat4からvec4(クォータニオン)への変換
vec4 mat4ToQuat(mat4 m) {
    float tr = m[0][0] + m[1][1] + m[2][2];
    float qw, qx, qy, qz;
    if(tr > 0) {
        float S = sqrt(tr + 1.0) * 2; // S=4*qw
        qw = 0.25 * S;
        qx = (m[2][1] - m[1][2]) / S;
        qy = (m[0][2] - m[2][0]) / S;
        qz = (m[1][0] - m[0][1]) / S;
    } else if((m[0][0] > m[1][1]) && (m[0][0] > m[2][2])) {
        float S = sqrt(1.0 + m[0][0] - m[1][1] - m[2][2]) * 2; // S=4*qx
        qw = (m[2][1] - m[1][2]) / S;
        qx = 0.25 * S;
        qy = (m[0][1] + m[1][0]) / S;
        qz = (m[0][2] + m[2][0]) / S;
    } else if(m[1][1] > m[2][2]) {
        float S = sqrt(1.0 + m[1][1] - m[0][0] - m[2][2]) * 2; // S=4*qy
        qw = (m[0][2] - m[2][0]) / S;
        qx = (m[0][1] + m[1][0]) / S;
        qy = 0.25 * S;
        qz = (m[1][2] + m[2][1]) / S;
    } else {
        float S = sqrt(1.0 + m[2][2] - m[0][0] - m[1][1]) * 2; // S=4*qz
        qw = (m[1][0] - m[0][1]) / S;
        qx = (m[0][2] + m[2][0]) / S;
        qy = (m[1][2] + m[2][1]) / S;
        qz = 0.25 * S;
    }
    return vec4(qx, qy, qz, qw);
}

// vec4(クォータニオン)からmat4への変換
mat4 quatToMat4(vec4 q) {
    float x = q.x;
    float y = q.y;
    float z = q.z;
    float w = q.w;
    return mat4(1.0 - 2.0 * y * y - 2.0 * z * z, 2.0 * x * y - 2.0 * z * w, 2.0 * x * z + 2.0 * y * w, 0.0, 2.0 * x * y + 2.0 * z * w, 1.0 - 2.0 * x * x - 2.0 * z * z, 2.0 * y * z - 2.0 * x * w, 0.0, 2.0 * x * z - 2.0 * y * w, 2.0 * y * z + 2.0 * x * w, 1.0 - 2.0 * x * x - 2.0 * y * y, 0.0, 0.0, 0.0, 0.0, 1.0);
}

// テクスチャからボーン変形行列を取得する
mat4 getBoneMatrix(int boneIndex) {
    int rowIndex = boneIndex * 4 / boneMatrixWidth;
    int colIndex = (boneIndex * 4) - (boneMatrixWidth * rowIndex);

    vec4 row0 = texelFetch(boneMatrixTexture, ivec2(colIndex + 0, rowIndex), 0);
    vec4 row1 = texelFetch(boneMatrixTexture, ivec2(colIndex + 1, rowIndex), 0);
    vec4 row2 = texelFetch(boneMatrixTexture, ivec2(colIndex + 2, rowIndex), 0);
    vec4 row3 = texelFetch(boneMatrixTexture, ivec2(colIndex + 3, rowIndex), 0);
    mat4 boneMatrix = mat4(row0, row1, row2, row3);

    return boneMatrix;
}

// 移動行列を作成する
mat4 createTranslationMatrix(vec3 translation) {
    mat4 translationMatrix = mat4(1.0);
    translationMatrix[3][0] = translation.x;
    translationMatrix[3][1] = translation.y;
    translationMatrix[3][2] = translation.z;
    return translationMatrix;
}

// スケール行列を作成する
mat4 createScaleMatrix(vec3 scale) {
    mat4 scaleMatrix = mat4(1.0);
    scaleMatrix[0][0] = scale.x;
    scaleMatrix[1][1] = scale.y;
    scaleMatrix[2][2] = scale.z;
    return scaleMatrix;
}

// スケールベクトルを計算する
vec3 calculateScaleVector(mat4 mat) {
    vec3 scale;
    scale.x = length(vec3(mat[0][0], mat[0][1], mat[0][2]));
    scale.y = length(vec3(mat[1][0], mat[1][1], mat[1][2]));
    scale.z = length(vec3(mat[2][0], mat[2][1], mat[2][2]));
    return scale;
}

// SDEF変形中心Cの計算
vec3 calculateSdefC(mat4 boneMatrix0, mat4 boneMatrix1, float boneWeight0, float boneWeight1, vec4 transformedPosition0, vec4 transformedPosition1) {

    // ボーンの位置を抽出
    vec3 c0 = (boneMatrix0 * vec4(sdefC, 1.0)).xyz * boneWeight0;
    vec3 c1 = (boneMatrix1 * vec4(sdefC, 1.0)).xyz * boneWeight1;

    // C点をボーンのウェイトに基づいて補間
    vec3 interpolatedC = (c0 + c1) / (boneWeight0 + boneWeight1);

    return interpolatedC;
}

// C点から見たR0とR1の補間を行い、C点の補正を適用する
vec3 interpolateCPoint(vec3 interpolatedC, mat4 boneMatrix0, mat4 boneMatrix1, float boneWeight0, float boneWeight1, vec4 transformedPosition0, vec4 transformedPosition1) {
    vec3 vecR0inB0 = sdefR0 - sdefB0;
    vec3 vecCinB0 = sdefC - sdefB0;
    vec3 vecR1inB1 = sdefR1 - sdefB1;
    vec3 vecCinB1 = sdefC - sdefB1;
    vec3 vecR0inC = sdefR0 - sdefC;
    vec3 vecR1inC = sdefR1 - sdefC;

    // R0/R1影響係数算出
    float len0 = length(vecR0inB0 - vecCinB0);
    float len1 = length(vecR1inB1 - vecCinB1);

    float r1Bias = 0.0;
    if(len0 > 0.0 && len1 == 0.0) {
        r1Bias = 1.0;
    } else if(len0 == 0.0 && len1 > 0.0) {
        r1Bias = 0.0;
    } else if(len0 + len1 != 0.0) {
        float bias = len0 / (len0 + len1);
        if(!isinf(bias) && !isnan(bias)) {
            r1Bias = clamp(bias, 0.0, 1.0);
        }
    }
    float r0Bias = 1.0 - r1Bias;

    // C点に基づいて変形されたR0とR1を計算
    vec3 transformedR0 = (boneMatrix0 * vec4(vecR0inC, 0.0)).xyz;
    vec3 transformedR1 = (boneMatrix1 * vec4(vecR1inC, 0.0)).xyz;

    // C点の補正：ウェイトに基づいてR0とR1の変形後の位置の平均を取る
    vec3 weightedAverage = (transformedR0 * r0Bias + transformedR1 * r1Bias) * 0.5;

    // 最終的なC点の位置：補正Cに変形後のウェイト付き平均を加算
    vec3 correctedC = interpolatedC + weightedAverage;

    return correctedC;
}

// クォータニオンによるボーンの回転を計算し、頂点Pを変形させる
vec4 applySdefRotation(vec3 correctedC, mat4 boneMatrix0, mat4 boneMatrix1, float boneWeight0, float boneWeight1) {
    vec3 vecPinC = position - sdefC;

    // ボーンのクォータニオン回転を取得
    vec4 boneQuat0 = mat4ToQuat(boneMatrix0);
    vec4 boneQuat1 = mat4ToQuat(boneMatrix1);

    // ボーンのウェイトに基づいてクォータニオンをSLERPにより補間
    vec4 slerpedQuat = slerp(boneQuat0, boneQuat1, boneWeight1);

    // クォータニオンを回転行列に変換
    mat4 rotationMatrix = quatToMat4(slerpedQuat);

    // 回転行列を使用して頂点を変形
    vec4 rotatedPosition = rotationMatrix * vec4(vecPinC, 1.0) + vec4(correctedC, 0.0);

    return rotatedPosition;
}

void main() {
    vec4 position4 = vec4(position, 1.0);

    // 各頂点で使用されるボーン変形行列を計算する
    totalBoneWeight = 0;
    mat4 boneTransformMatrix = mat4(0.0);
    mat3 normalTransformMatrix = mat3(1.0);

    if(isSdef == 1.0) {
        // SDEFの場合は、SDEF用の頂点位置を計算する

        // ボーンインデックスからボーン変形行列を取得
        mat4 boneMatrix0 = getBoneMatrix(int(boneIndexes[0]));
        mat4 boneMatrix1 = getBoneMatrix(int(boneIndexes[1]));

        float boneWeight0 = boneWeights[0];
        float boneWeight1 = boneWeights[1];

        vec4 transformedPosition0 = boneMatrix0 * vec4(0.0, 0.0, 0.0, 1.0);
        vec4 transformedPosition1 = boneMatrix1 * vec4(0.0, 0.0, 0.0, 1.0);

        // 1. 頂点Pに対する変形中心Cを計算
        vec3 interpolatedC = calculateSdefC(boneMatrix0, boneMatrix1, boneWeight0, boneWeight1, transformedPosition0, transformedPosition1);

        // 変形中心Cの補正を計算
        vec3 correctedC = interpolateCPoint(interpolatedC, boneMatrix0, boneMatrix1, boneWeight0, boneWeight1, transformedPosition0, transformedPosition1);

        // ボーンの回転を適用して頂点Pを変形させる
        vec4 vecPosition = applySdefRotation(correctedC, boneMatrix0, boneMatrix1, boneWeight0, boneWeight1);

        gl_Position = modelViewProjectionMatrix * modelViewMatrix * vecPosition;
    } else {
        for(int i = 0; i < 4; i++) {
            float boneWeight = boneWeights[i];
            int boneIndex = int(boneIndexes[i]);

            // テクスチャからボーン変形行列を取得する
            mat4 boneMatrix = getBoneMatrix(boneIndex);

            // ボーン変形行列を加算する
            boneTransformMatrix += boneMatrix * boneWeight;
        }

        gl_Position = modelViewProjectionMatrix * modelViewMatrix * boneTransformMatrix * position4;

        // 各頂点で使用される法線変形行列をボーン変形行列から回転情報のみ抽出して生成する
        normalTransformMatrix = mat3(boneTransformMatrix);
    }

    // 頂点法線
    vetexNormal = normalize(normalTransformMatrix * normalize(normal)).xyz;

    // 頂点色設定(透過込み)
    vertexColor = clamp(diffuse, 0.0, 1.0);

    if(0 == useToon) {
        // ディフューズ色＋アンビエント色 計算
        float lightNormal = clamp(dot(vetexNormal, -lightDirection), 0.0, 1.0);
        vertexColor.rgb += diffuse.rgb * lightNormal;
        vertexColor = clamp(vertexColor, 0.0, 1.0);
    }

    // テクスチャ描画位置
    vertexUv = uv;

    if(1 == useSphere) {
        // Sphereマップ計算
        if(3 == sphereMode) {
            // PMXサブテクスチャ座標
            sphereUv = extendUv;
        } else {
	        // スフィアマップテクスチャ座標
            vec3 normalWv = mat3(modelViewMatrix) * vetexNormal;
            sphereUv.x = normalWv.x * 0.5 + 0.5;
            sphereUv.y = normalWv.y * -0.5 + 0.5;
        }
        // sphereUv += morphUv1.xy;
    }

    // カメラとの相対位置
    vec3 eye = cameraPosition - (boneTransformMatrix * position4).xyz;

    // スペキュラ色計算
    vec3 HalfVector = normalize(normalize(eye) + -lightDirection);
    vertexSpecular = pow(max(0, dot(HalfVector, vetexNormal)), max(0.000001, specular.w)) * specular.rgb;
}
