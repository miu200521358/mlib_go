#version 440 core

uniform mat4 modelViewProjectionMatrix;
uniform mat4 modelViewMatrix;

// ボーン変形行列を格納するテクスチャ
uniform sampler2D boneMatrixTexture;
uniform int boneMatrixWidth;
uniform int boneMatrixHeight;

in layout(location = 0) vec3 position;
in layout(location = 1) vec3 normal;
in layout(location = 2) vec2  uv;
in layout(location = 3) vec2  extendUv;
in layout(location = 4) float vertexEdge;
in layout(location = 5) vec4  boneIndexes;
in layout(location = 6) vec4  boneWeights;

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

void main() {
    vec4 position4 = vec4(position, 1.0);

    // 各頂点で使用されるボーン変形行列を計算する
    totalBoneWeight = 0;
    mat4 boneTransformMatrix = mat4(0.0);
    for(int i = 0; i < 4; i++) {
        float boneWeight = boneWeights[i];
        int boneIndex = int(boneIndexes[i]);

        // テクスチャからボーン変形行列を取得する
        int rowIndex = boneIndex * 4 / boneMatrixWidth;
        int colIndex = (boneIndex * 4) - (boneMatrixWidth * rowIndex);

        vec4 row0 = texelFetch(boneMatrixTexture, ivec2(colIndex + 0, rowIndex), 0);
        vec4 row1 = texelFetch(boneMatrixTexture, ivec2(colIndex + 1, rowIndex), 0);
        vec4 row2 = texelFetch(boneMatrixTexture, ivec2(colIndex + 2, rowIndex), 0);
        vec4 row3 = texelFetch(boneMatrixTexture, ivec2(colIndex + 3, rowIndex), 0);
        mat4 boneMatrix = mat4(row0, row1, row2, row3);

        // ボーン変形行列を乗算する
        boneTransformMatrix += boneMatrix * boneWeight;
    }

    gl_Position = modelViewProjectionMatrix * modelViewMatrix * boneTransformMatrix * position4;

    // 各頂点で使用される法線変形行列をボーン変形行列から回転情報のみ抽出して生成する
    mat3 normalTransformMatrix = mat3(boneTransformMatrix);

    // 頂点法線
    vetexNormal = normalize(normalTransformMatrix * normalize(normal)).xyz;

    // 頂点色設定(透過込み)
    vertexColor = clamp(diffuse, 0.0, 1.0);

    if (0 == useToon) {
        // ディフューズ色＋アンビエント色 計算
        float lightNormal = clamp(dot( vetexNormal, -lightDirection ), 0.0, 1.0);
        vertexColor.rgb += diffuse.rgb * lightNormal;
        vertexColor = clamp(vertexColor, 0.0, 1.0);
    }

    // テクスチャ描画位置
    vertexUv = uv;

    if (1 == useSphere) {
        // Sphereマップ計算
        if (3 == sphereMode) {
            // PMXサブテクスチャ座標
            sphereUv = extendUv;
        }
        else {
	        // スフィアマップテクスチャ座標
            vec3 normalWv = mat3(modelViewMatrix) * vetexNormal;
	        sphereUv.x = normalWv.x * 0.5f + 0.5f;
	        sphereUv.y = normalWv.y * -0.5f + 0.5f;
        }
        // sphereUv += morphUv1.xy;
    }

    // カメラとの相対位置
    vec3 eye = cameraPosition - (boneTransformMatrix * position4).xyz;

    // スペキュラ色計算
    vec3 HalfVector = normalize( normalize(eye) + -lightDirection );
    vertexSpecular = pow( max(0, dot( HalfVector, vetexNormal )), max(0.000001, specular.w) ) * specular.rgb;
}
