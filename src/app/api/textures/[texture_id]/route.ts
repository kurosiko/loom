import { NextResponse } from "next/server";
import type { Texture } from "@/types";

// Mock database
const textures: Texture[] = [
  {
    id: "1",
    name: "General Discussion",
    description: "A place for general conversations",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: "2",
    name: "Tech Talk",
    description: "Discuss technology and programming",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
  {
    id: "3",
    name: "Creative Corner",
    description: "Share your creative projects",
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  },
];

export async function GET(
  request: Request,
  { params }: { params: Promise<{ texture_id: string }> }
) {
  const { texture_id } = await params;
  const texture = textures.find((t) => t.id === texture_id);

  if (!texture) {
    return NextResponse.json({ error: "Texture not found" }, { status: 404 });
  }

  return NextResponse.json(texture);
}

export async function PUT(
  request: Request,
  { params }: { params: Promise<{ texture_id: string }> }
) {
  const { texture_id } = await params;
  const body = await request.json();
  const index = textures.findIndex((t) => t.id === texture_id);

  if (index === -1) {
    return NextResponse.json({ error: "Texture not found" }, { status: 404 });
  }

  textures[index] = {
    ...textures[index],
    ...body,
    updatedAt: new Date().toISOString(),
  };

  return NextResponse.json(textures[index]);
}

export async function DELETE(
  request: Request,
  { params }: { params: Promise<{ texture_id: string }> }
) {
  const { texture_id } = await params;
  const index = textures.findIndex((t) => t.id === texture_id);

  if (index === -1) {
    return NextResponse.json({ error: "Texture not found" }, { status: 404 });
  }

  textures.splice(index, 1);
  return NextResponse.json({ success: true });
}
