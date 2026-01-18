#version 430 core

out vec4 OutColor;

in vec2 TexCoord;

uniform sampler2D overrideTexture;

void main() {
    vec4 texColor = texture(overrideTexture, TexCoord);
    OutColor = vec4(texColor.rgb, 0.4);
}