#version 440 core

in vec4 boneColor;
out vec4 outColor;

void main() {
    if(boneColor.a < 1e-6) {
        discard;
    }

    outColor = boneColor;
}