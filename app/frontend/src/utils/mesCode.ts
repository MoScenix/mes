export type MesCodeKind = 'FLOW' | 'ITEM_UNIT' | 'ENGINEERING_ORDER'
export type MesDetailKind = MesCodeKind | 'WORK_ORDER' | 'ITEM' | 'PROCESS'

export type ParsedMesCode = {
  type: 'flow' | 'itemUnit' | 'engineeringOrder' | 'numeric' | 'unknown'
  kind?: MesCodeKind
  id: number
}

const typeByKind: Record<MesCodeKind, ParsedMesCode['type']> = {
  FLOW: 'flow',
  ITEM_UNIT: 'itemUnit',
  ENGINEERING_ORDER: 'engineeringOrder',
}

export const makeMesCode = (kind: MesCodeKind, id?: number) => (id ? `MES:${kind}:${id}` : '')

export const makeItemUnitCode = (id: number) => makeMesCode('ITEM_UNIT', id)

export const makeFlowCode = (id: number) => makeMesCode('FLOW', id)

export const makeEngineeringOrderCode = (id: number) => makeMesCode('ENGINEERING_ORDER', id)

export const parseMesCode = (value: string, defaultKind?: MesCodeKind): ParsedMesCode => {
  const raw = value.trim().replace(/^['"]|['"]$/g, '')
  let text = raw
  try {
    text = decodeURIComponent(raw)
  } catch {
    text = raw
  }
  const matched = text.match(/MES:(FLOW|ITEM_UNIT|ENGINEERING_ORDER):(\d+)/i)
  if (matched) {
    const kind = matched[1].toUpperCase() as MesCodeKind
    return { type: typeByKind[kind], kind, id: Number(matched[2]) }
  }

  try {
    const json = JSON.parse(text) as { kind?: string; type?: string; id?: number | string }
    const candidate = String(json.kind || json.type || '').toUpperCase() as MesCodeKind
    const id = Number(json.id)
    if (typeByKind[candidate] && id > 0) return { type: typeByKind[candidate], kind: candidate, id }
  } catch {
    // Not a JSON code payload.
  }

  const queryKind = text.match(/[?&#](?:kind|type)=(FLOW|ITEM_UNIT|ENGINEERING_ORDER)/i)?.[1]
  const queryId = Number(text.match(/[?&#]id=(\d+)/i)?.[1])
  if (queryKind && queryId > 0) {
    const kind = queryKind.toUpperCase() as MesCodeKind
    return { type: typeByKind[kind], kind, id: queryId }
  }

  const numericId = Number(text)
  if (Number.isFinite(numericId) && numericId > 0) {
    return {
      type: defaultKind ? typeByKind[defaultKind] : 'numeric',
      kind: defaultKind,
      id: numericId,
    }
  }

  return { type: 'unknown', id: 0 }
}

export const parseMesId = (value: string, defaultKind?: MesCodeKind) => {
  const parsed = parseMesCode(value, defaultKind)
  return parsed.id || undefined
}
