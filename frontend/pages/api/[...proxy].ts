import type { NextApiRequest, NextApiResponse } from "next";
import httpProxy from "http-proxy";

export const config = {
  api: {
    // Enable `externalResolver` option in Next.js
    externalResolver: true,
    bodyParser: false,
  },
};

export default async function proxy(req: NextApiRequest, res: NextApiResponse) {
  try {
    const nextEncryptedToken = req.cookies['next-auth.session-token']
    req.headers[`Authorization`] = `Bearer ${nextEncryptedToken}`
    req.headers['content-type'] = `application/json`
    
    return new Promise((resolve, reject) => {
      const proxy: httpProxy = httpProxy.createProxy();
      proxy.once("proxyRes", resolve).once("error", reject).web(req, res, {
        changeOrigin: true,
        target: process.env.EXTERNAL_API_ENDPOINT,
      });
    });
  } catch(err) {
    console.error({ err })
    return res.status(500).json({ error: `server error`})
  }
}