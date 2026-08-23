#version 430 core

uniform mat4 projectionMatrix;
uniform mat4 viewMatrix;
uniform vec2 viewportSize;

in layout(location = 0) vec3 startPosition;
in layout(location = 1) vec3 endPosition;
in layout(location = 2) vec4 startColor;
in layout(location = 3) vec4 endColor;
in layout(location = 4) vec2 corner;
in layout(location = 5) float sizePixels;
in layout(location = 6) float markerMode;

out vec4 color4;
out vec2 markerCoordinate;
flat out float isMarker;

void main() {
    vec4 startClip = projectionMatrix * viewMatrix * vec4(startPosition, 1.0);
    vec4 endClip = projectionMatrix * viewMatrix * vec4(endPosition, 1.0);
    isMarker = markerMode;
    markerCoordinate = vec2(0.0);

    if (markerMode > 0.5) {
        vec2 offset = corner * sizePixels / viewportSize;
        gl_Position = startClip;
        gl_Position.xy += offset * startClip.w;
        color4 = startColor;
        markerCoordinate = corner;
        return;
    }

    vec2 startNdc = startClip.xy / startClip.w;
    vec2 endNdc = endClip.xy / endClip.w;
    vec2 direction = endNdc - startNdc;
    vec2 normal = vec2(0.0, 1.0);
    if (dot(direction, direction) > 1e-12) {
        normal = normalize(vec2(-direction.y, direction.x));
    }
    vec4 baseClip = mix(startClip, endClip, corner.x);
    vec2 offset = normal * corner.y * sizePixels / viewportSize;
    baseClip.xy += offset * baseClip.w;
    gl_Position = baseClip;
    color4 = mix(startColor, endColor, corner.x);
}
