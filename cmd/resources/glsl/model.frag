#version 440

uniform vec4 diffuse;
uniform vec3 ambient;

in float alpha;
in vec4 vertexColor;
in vec3 vertexSpecular;
in vec2 vertexUv;
in vec3 vetexNormal;
in vec2 sphereUv;
in vec3 eye;
in float totalBoneWeight;

out vec4  outColor;

void main() {
    outColor = vertexColor;

    // スペキュラ適用
    outColor.rgb += vertexSpecular;
}