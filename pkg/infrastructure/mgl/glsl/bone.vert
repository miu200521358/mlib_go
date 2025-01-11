#version 440 core

uniform mat4 projectionMatrix;
uniform mat4 viewMatrix;

// ボーン変形行列を格納するテクスチャ
uniform sampler2D boneMatrixTexture;
uniform int boneMatrixWidth;
uniform int boneMatrixHeight;

in layout(location = 0) vec3 position;
in layout(location = 1) vec4 boneIndexes;
in layout(location = 2) vec4 boneWeights;
in layout(location = 3) vec4 color;

out float totalBoneWeight;
out vec4 boneColor;

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

void main() {
    vec4 position4 = vec4(position, 1.0);

    // 各頂点で使用されるボーン変形行列を計算する
    totalBoneWeight = 0;
    mat4 boneTransformMatrix = mat4(0.0);

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

    // ボーンカラーを出力
    boneColor = color;
}
