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
  const text = value.trim()
  const matched = text.match(/^MES:(FLOW|ITEM_UNIT|ENGINEERING_ORDER):(\d+)$/i)
  if (matched) {
    const kind = matched[1].toUpperCase() as MesCodeKind
    return { type: typeByKind[kind], kind, id: Number(matched[2]) }
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
