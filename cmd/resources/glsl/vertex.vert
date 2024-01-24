#version 440

in vec3 position;
// in layout(location = 1) vec3  normal;
// in layout(location = 2) vec4  boneIndexes;
// in layout(location = 3) vec4  boneWeights;

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

uniform vec4 diffuse;

out vec4 vertexColor;

void main() {
    // 頂点色設定
    vertexColor = clamp(diffuse, 0.0, 1.0);

    gl_Position = projection * camera * model * vec4(position, 1);
}
