#version 440 core

uniform mat4 modelViewProjectionMatrix;
uniform mat4 modelViewMatrix;

// ボーン変形行列を格納するテクスチャ
uniform sampler2D boneMatrixTexture;
uniform int boneMatrixWidth;
uniform int boneMatrixHeight;

in layout(location = 0) float boneIndex;
in layout(location = 1) vec4 typeColor;
in layout(location = 2) vec4 rotation;
in layout(location = 3) vec3 position;

out vec4 rigidbodyColor;

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

void main() {
    mat4 boneTransformMatrix = mat4(0.0);
    if (boneIndex >= 0.0) {
        // テクスチャからボーン変形行列を取得する
        mat4 boneMatrix = getBoneMatrix(int(boneIndex));

        // ボーン変形行列を加算する
        boneTransformMatrix += boneMatrix;

        // ボーン変形行列を適用して描画
        gl_Position = modelViewProjectionMatrix * modelViewMatrix * boneTransformMatrix * quatToMat4(rotation) * vec4(position, 1.0);
    } else {
        // ボーンに紐付いていない場合、そのまま描画
        gl_Position = modelViewProjectionMatrix * modelViewMatrix * quatToMat4(rotation) * vec4(position, 1.0);
    }

    rigidbodyColor = typeColor;
}
