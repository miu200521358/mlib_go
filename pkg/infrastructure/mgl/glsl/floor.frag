#version 440 core

in vec4 color4;
out vec4  outColor;

uniform float windowOpacity;

void main() {
    outColor = color4;
    outColor *= windowOpacity;
}