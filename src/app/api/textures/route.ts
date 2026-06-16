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

export async function GET() {
  return NextResponse.json(textures);
}

export async function POST(request: Request) {
  const body = await request.json();
  const { name, description } = body;

  const newTexture: Texture = {
    id: String(textures.length + 1),
    name,
    description,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  textures.push(newTexture);
  return NextResponse.json(newTexture, { status: 201 });
}
