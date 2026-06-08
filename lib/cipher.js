export const DEFAULT_KEY = "instahost";
export const CHECKSUM_BYTES = 4;

const CRC32_TABLE = new Uint32Array(256);
for (let i = 0; i < 256; i++) {
  let c = i;
  for (let j = 0; j < 8; j++) {
    c = (c & 1) ? (0xedb88320 ^ (c >>> 1)) : (c >>> 1);
  }
  CRC32_TABLE[i] = c;
}

export function crc32(data) {
  const input = Buffer.isBuffer(data) ? data : Buffer.from(data);
  let crc = 0xffffffff;

  for (let i = 0; i < input.length; i++) {
    crc = CRC32_TABLE[(crc ^ input[i]) & 0xff] ^ (crc >>> 8);
  }

  return (crc ^ 0xffffffff) >>> 0;
}

export function sealWithChecksum(data) {
  const input = Buffer.isBuffer(data) ? data : Buffer.from(data);
  const out = Buffer.alloc(input.length + CHECKSUM_BYTES);
  input.copy(out);
  out.writeUInt32BE(crc32(input), input.length);
  return out;
}

export function openWithChecksum(data) {
  const input = Buffer.isBuffer(data) ? data : Buffer.from(data);
  if (input.length <= CHECKSUM_BYTES) {
    return null;
  }

  const payload = input.subarray(0, input.length - CHECKSUM_BYTES);
  const stored = input.readUInt32BE(input.length - CHECKSUM_BYTES);
  if (stored !== crc32(payload)) {
    return null;
  }

  return payload;
}

export function xorBytes(data, key = DEFAULT_KEY) {
  const keyBytes = Buffer.from(key, "utf8");
  const input = Buffer.isBuffer(data) ? data : Buffer.from(data);
  const out = Buffer.alloc(input.length);

  for (let i = 0; i < input.length; i++) {
    out[i] = input[i] ^ keyBytes[i % keyBytes.length];
  }

  return out;
}
