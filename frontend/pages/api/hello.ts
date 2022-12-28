// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next'
import { encode, decode, JWT } from 'next-auth/jwt'
import { SignJWT, JWTPayload } from "jose"
import { v4 as uuid } from "uuid"
import hkdf from "@panva/hkdf"

type Data = {
  name: string
}

const DEFAULT_MAX_AGE = 30 * 24 * 60 * 60 // 30 days

const now = () => (Date.now() / 1000) | 0

/* next-auth key derivation function */
async function getNextAuthEncryptionKey(secret: string | Buffer) {
  return await hkdf(
    "sha256",
    secret,
    "",
    "NextAuth.js Generated Encryption Key",
    32
  )
}

async function signExternalJWT(jwt: JWT, secret: Uint8Array) {
  return await new SignJWT(jwt as JWTPayload)
  .setProtectedHeader({ alg: 'HS256'})
  .setIssuedAt()
  .setExpirationTime(now() + DEFAULT_MAX_AGE)
  .setJti(uuid())
  .sign(secret)
}


export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Data>
) {
  const nextEncryptedToken = req.cookies['next-auth.session-token']

  const result = await fetch(`${process.env.EXTERNAL_API_ENDPOINT}/user` as string, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${nextEncryptedToken}`
    }
  })
  try {
    const data = await result.json();
    console.log({ data })
    res.status(200).json({ goApiResponse: data } as any)
  } catch(err) {
    res.status(200).json({ data: err } as any)
  }
}
