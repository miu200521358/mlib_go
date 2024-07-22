#version 440 core

uniform vec4 edgeColor;
out vec4 outColor;

uniform float windowOpacity;

void main() {
    if(edgeColor.a < 1e-6) {
        discard;
    }

    outColor = edgeColor;
    outColor.a *= windowOpacity;
}