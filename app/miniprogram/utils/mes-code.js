const KINDS = ['FLOW', 'ITEM_UNIT', 'ENGINEERING_ORDER']

function parseMesCode(value, defaultKind) {
  const raw = String(value || '').trim().replace(/^['"]|['"]$/g, '')
  let text = raw
  try { text = decodeURIComponent(raw) } catch (error) {}
  const matched = text.match(/MES:(FLOW|ITEM_UNIT|ENGINEERING_ORDER):(\d+)/i)
  if (matched) return { kind: matched[1].toUpperCase(), id: Number(matched[2]) }
  try {
    const json = JSON.parse(text)
    const kind = String(json.kind || json.type || '').toUpperCase()
    const id = Number(json.id)
    if (KINDS.includes(kind) && id > 0) return { kind, id }
  } catch (error) {}
  const kindMatch = text.match(/[?&#](?:kind|type)=(FLOW|ITEM_UNIT|ENGINEERING_ORDER)/i)
  const idMatch = text.match(/[?&#]id=(\d+)/i)
  if (kindMatch && idMatch) return { kind: kindMatch[1].toUpperCase(), id: Number(idMatch[1]) }
  const id = Number(text)
  if (defaultKind && Number.isFinite(id) && id > 0) return { kind: defaultKind, id }
  return { kind: '', id: 0 }
}

function scan(expectedKind) {
  return new Promise((resolve, reject) => {
    wx.scanCode({
      scanType: ['qrCode'],
      success(result) {
        const parsed = parseMesCode(result.result, expectedKind)
        if (!parsed.id || (expectedKind && parsed.kind !== expectedKind)) {
          reject(new Error('二维码类型不正确'))
          return
        }
        resolve(parsed)
      },
      fail(error) {
        if (error.errMsg && error.errMsg.includes('cancel')) {
          reject(null)
          return
        }
        reject(new Error('扫码失败，请重试'))
      }
    })
  })
}

module.exports = { KINDS, parseMesCode, scan }
