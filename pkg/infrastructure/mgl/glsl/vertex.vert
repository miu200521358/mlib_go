#version 440 core

uniform mat4 projectionMatrix;
uniform mat4 viewMatrix;

// ボーン変形行列を格納するテクスチャ
uniform sampler2D boneMatrixTexture;
uniform int boneMatrixWidth;
uniform int boneMatrixHeight;

uniform vec3 cursorPositions[50];

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
in layout(location = 11) vec3 vertexDelta;  // 頂点モーフ
in layout(location = 12) vec4 uvDelta; // UVモーフ
in layout(location = 13) vec4 uv1Delta; // 拡張UV1モーフ
in layout(location = 14) vec3 afterVertexDelta; // ボーン変形後頂点モーフ

layout(std430, binding = 0) buffer VertexBuffer {
    vec4 vertexPositions[];
};

out float totalBoneWeight;
out float vertexUvX;

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

// クォータニオンによるボーンの回転を計算し、頂点Pを変形させる
mat4 calculateSdefMatrix(mat4 boneMatrix0, mat4 boneMatrix1, float boneWeight0, float boneWeight1) {
    // ボーンのクォータニオン回転を取得
    vec4 boneQuat0 = mat4ToQuat(boneMatrix0);
    vec4 boneQuat1 = mat4ToQuat(boneMatrix1);

    // ボーンのウェイトに基づいてクォータニオンをSLERPにより補間
    vec4 slerpedQuat1 = slerp(boneQuat0, boneQuat1, boneWeight1);

    // クォータニオンを回転行列に変換
    mat4 rotationMatrix = quatToMat4(slerpedQuat1);

    return rotationMatrix;
}

// 補間点R0/R1をBDEF2移動させて交点Cを補正する
vec4 calculateCorrectedC(mat4 boneMatrix0, mat4 boneMatrix1, float boneWeight0, float boneWeight1) {
    // R0/R1影響係数算出
    float lenR0C = length(sdefR0 - sdefC);
    float lenR1C = length(sdefR1 - sdefC);

    float r1Bias = 0.0;
    if(lenR1C == 0.0) {
        r1Bias = 0.0;
    } else if(lenR0C == 0.0) {
        r1Bias = 1.0;
    } else if(lenR0C + lenR1C != 0.0) {
        float bias = lenR0C / (lenR0C + lenR1C);
        if(!isinf(bias) && !isnan(bias)) {
            r1Bias = clamp(bias, 0.0, 1.0);
        }
    }
    float r0Bias = 1.0 - r1Bias;

    // ウェイトに基づいたC (BDEF2移動させたC)
    vec4 weightedC0 = boneMatrix0 * boneWeight0 * vec4(sdefC, 1.0);
    vec4 weightedC1 = boneMatrix1 * boneWeight1 * vec4(sdefC, 1.0);
    vec4 weightedC = weightedC0 + weightedC1;

    // 影響係数に基づいたR
    vec4 biasR0 = boneMatrix0 * r0Bias * vec4(sdefR0, 1.0);
    vec4 biasR1 = boneMatrix1 * r1Bias * vec4(sdefR1, 1.0);
    vec4 biasR = biasR0 + biasR1;

    // return biasR;
    return (weightedC + biasR) * 0.5;
}

// 逆行列
// 4x4行列の逆行列を求める関数
mat4 inverse(mat4 m) {
    float
    a00 = m[0][0], a01 = m[0][1], a02 = m[0][2], a03 = m[0][3],
    a10 = m[1][0], a11 = m[1][1], a12 = m[1][2], a13 = m[1][3],
    a20 = m[2][0], a21 = m[2][1], a22 = m[2][2], a23 = m[2][3],
    a30 = m[3][0], a31 = m[3][1], a32 = m[3][2], a33 = m[3][3];

    float
    b00 = a00 * a11 - a01 * a10,
    b01 = a00 * a12 - a02 * a10,
    b02 = a00 * a13 - a03 * a10,
    b03 = a01 * a12 - a02 * a11,
    b04 = a01 * a13 - a03 * a11,
    b05 = a02 * a13 - a03 * a12,
    b06 = a20 * a31 - a21 * a30,
    b07 = a20 * a32 - a22 * a30,
    b08 = a20 * a33 - a23 * a30,
    b09 = a21 * a32 - a22 * a31,
    b10 = a21 * a33 - a23 * a31,
    b11 = a22 * a33 - a23 * a32;

    float
    det = b00 * b11 - b01 * b10 + b02 * b09 + b03 * b08 - b04 * b07 + b05 * b06;

    return mat4(
        a11 * b11 - a12 * b10 + a13 * b09,
        a02 * b10 - a01 * b11 - a03 * b09,
        a31 * b05 - a32 * b04 + a33 * b03,
        a22 * b04 - a21 * b05 - a23 * b03,
        a12 * b08 - a10 * b11 - a13 * b07,
        a00 * b11 - a02 * b08 + a03 * b07,
        a32 * b02 - a30 * b05 - a33 * b01,
        a20 * b05 - a22 * b02 + a23 * b01,
        a10 * b10 - a11 * b08 + a13 * b06,
        a01 * b08 - a00 * b10 - a03 * b06,
        a30 * b04 - a31 * b02 + a33 * b00,
        a21 * b02 - a20 * b04 - a23 * b00,
        a11 * b07 - a10 * b09 - a12 * b06,
        a00 * b09 - a01 * b07 + a02 * b06,
        a31 * b01 - a30 * b03 - a32 * b00,
        a20 * b03 - a21 * b01 + a22 * b00) / det;
}

float distanceToVectors(vec3 point) {
    float minDistance = 1000000.0;
    for(int i = 0; i < 50; i++) {
        vec3 cursorPosition = cursorPositions[i-1];
        if (cursorPosition.x == 0 && cursorPosition.y == 0 && cursorPosition.z == 0) {
            continue;
        }
        float pointDistance = distance(point, cursorPosition);
        if(pointDistance < minDistance) {
            minDistance = pointDistance;
        }
    }
    return minDistance;
}

void main() {
    vec4 position4 = vec4(position + vertexDelta, 1.0);

    // 各頂点で使用されるボーン変形行列を計算する
    totalBoneWeight = 0;
    mat4 boneTransformMatrix = mat4(0.0);
    mat3 normalTransformMatrix = mat3(1.0);

    vec4 vecGlobalPosition;

    if(isSdef == 1.0) {
        // SDEFの場合は、SDEF用の頂点位置を計算する

        // ボーンインデックスからボーン変形行列を取得
        mat4 boneMatrix0 = getBoneMatrix(int(boneIndexes[0]));
        mat4 boneMatrix1 = getBoneMatrix(int(boneIndexes[1]));

        float boneWeight0 = boneWeights[0];
        float boneWeight1 = boneWeights[1];

        // ボーンの回転を適用して頂点Pを変形させる
        mat4 rotationMatrix = calculateSdefMatrix(boneMatrix0, boneMatrix1, boneWeight0, boneWeight1);

        // 補正Cを求める
        vec4 correctedC = calculateCorrectedC(boneMatrix0, boneMatrix1, boneWeight0, boneWeight1);

        // 回転行列を使用して頂点を変形
        vec4 rotatedPosition = rotationMatrix * vec4(position, 1.0);
        vec4 rotatedC = rotationMatrix * vec4(sdefC, 1.0);

        vec4 vecPosition = rotatedPosition - rotatedC + correctedC;

        // 頂点位置
        gl_Position = projectionMatrix * viewMatrix * vec4(vecPosition.xyz, 1.0);

        // SDEF結果頂点位置を格納
        vecGlobalPosition = vec4(vecPosition.xyz, 1.0);
    } else {
        for(int i = 0; i < 4; i++) {
            float boneWeight = boneWeights[i];
            int boneIndex = int(boneIndexes[i]);

            // テクスチャからボーン変形行列を取得する
            mat4 boneMatrix = getBoneMatrix(boneIndex);

            // ボーン変形行列を加算する
            boneTransformMatrix += boneMatrix * boneWeight;
        }

        // 頂点位置
        gl_Position = projectionMatrix * viewMatrix * boneTransformMatrix * position4;

        // ボーン変形行列を加味した頂点位置を格納
        vecGlobalPosition = boneTransformMatrix * position4;
    }

    // カーソルの矩形範囲内にある場合のみ、頂点位置を格納する
    float nearestDistance = distanceToVectors(vecGlobalPosition.xyz);
    if (nearestDistance < 0.1) {
        vertexPositions[gl_VertexID] = vec4(vecGlobalPosition.xyz, nearestDistance);
    } else {
        vertexPositions[gl_VertexID] = vec4(-1);
    }

    vertexUvX = uv.x + uvDelta.x;
}
