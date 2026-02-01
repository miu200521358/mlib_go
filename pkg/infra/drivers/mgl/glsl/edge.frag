#version 430 core

uniform vec4 edgeColor;
out vec4 outColor;

void main() {
    if(edgeColor.a < 1e-6) {
        discard;
    }

    outColor = edgeColor;
}