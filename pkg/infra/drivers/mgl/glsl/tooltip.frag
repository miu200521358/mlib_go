#version 430 core

in vec2 vTexCoord;
out vec4 outColor;

uniform sampler2D tooltipTexture;

void main() {
    vec4 color = texture(tooltipTexture, vTexCoord);
    if (color.a <= 0.0) {
        discard;
    }
    outColor = color;
}

