#version 440 core

uniform vec4 diffuse;
uniform vec3 ambient;

uniform int useTexture;
uniform sampler2D textureSampler;

uniform int useToon;
uniform sampler2D toonSampler;

uniform int useSphere;
uniform int sphereMode;
uniform sampler2D sphereSampler;

uniform vec3 lightDirection;

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

    if (1 == useTexture) {
        // テクスチャ適用
        outColor *= texture(textureSampler, vertexUv);
    }

    if (1 == useSphere) {
        // Sphere適用
        vec4 texColor = texture(sphereSampler, sphereUv);
        if (2 == sphereMode) {
            // スフィア加算
            outColor.rgb += texColor.rgb;
        }
        else {
            // スフィア乗算
            outColor.rgb *= texColor.rgb;
        }
        outColor.a *= texColor.a;
    }

    if (1 == useToon) {
        // Toon適用
        float lightNormal = dot( vetexNormal, -lightDirection );
        outColor *= texture(toonSampler, vec2(0, lightNormal));
    }

    // スペキュラ適用
    outColor.rgb += vertexSpecular;
}