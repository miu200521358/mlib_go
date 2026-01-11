package model

import "github.com/miu200521358/mlib_go/pkg/domain/model/collection"

// VertexCollection は頂点の IndexedCollection を表す。
type VertexCollection = collection.IndexedCollection[*Vertex]

// FaceCollection は面の IndexedCollection を表す。
type FaceCollection = collection.IndexedCollection[*Face]

// TextureCollection はテクスチャの NamedCollection を表す。
type TextureCollection = collection.NamedCollection[*Texture]

// MaterialCollection は材質の NamedCollection を表す。
type MaterialCollection = collection.NamedCollection[*Material]

// MorphCollection はモーフの NamedCollection を表す。
type MorphCollection = collection.NamedCollection[*Morph]

// DisplaySlotCollection は表示枠の NamedCollection を表す。
type DisplaySlotCollection = collection.NamedCollection[*DisplaySlot]

// RigidBodyCollection は剛体の NamedCollection を表す。
type RigidBodyCollection = collection.NamedCollection[*RigidBody]

// JointCollection はジョイントの NamedCollection を表す。
type JointCollection = collection.NamedCollection[*Joint]
