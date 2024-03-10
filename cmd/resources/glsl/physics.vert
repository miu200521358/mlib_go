#version 440 core

uniform mat4 modelViewProjectionMatrix;
uniform mat4 modelViewMatrix;

in layout(location = 0) vec3 position;

in vec3 color;
in float alpha;

out vec4 physicsColor;

void main() {
    gl_Position = modelViewProjectionMatrix * modelViewMatrix * vec4(position, 1.0);

    physicsColor = vec4(color, alpha);
}
