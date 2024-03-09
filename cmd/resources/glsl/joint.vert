#version 440 core

uniform mat4 modelViewProjectionMatrix;
uniform mat4 modelViewMatrix;

in layout(location = 0) vec4 typeColor;
in layout(location = 1) vec3 position;

out vec4 rigidbodyColor;

void main() {
    gl_Position = modelViewProjectionMatrix * modelViewMatrix * vec4(position, 1.0);

    rigidbodyColor = typeColor;
}
