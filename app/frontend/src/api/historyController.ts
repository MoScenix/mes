// @ts-ignore
/* eslint-disable */
import request from '@/request'

/** POST /history/add */
export async function addHistory(body: { initPrompt?: string }, options?: { [key: string]: any }) {
  return request<API.BaseResponseLong>('/history/add', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    data: body,
    ...(options || {}),
  })
}

/** POST /history/delete */
export async function deleteHistory(body: API.DeleteRequest, options?: { [key: string]: any }) {
  return request<API.BaseResponseBoolean>('/history/delete', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    data: body,
    ...(options || {}),
  })
}

/** POST /history/my/list/page/vo */
export async function listMyHistoryVoByPage(
  body: { pageNum?: number; pageSize?: number; historyName?: string },
  options?: { [key: string]: any },
) {
  return request<API.BaseResponsePageHistoryVO>('/history/my/list/page/vo', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    data: body,
    ...(options || {}),
  })
}

/** POST /history/admin/list/page/vo */
export async function listAllHistoryMessagesByPageForAdmin(
  body: API.HistoryMessageQueryRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponsePageHistoryMessage>('/history/admin/list/page/vo', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    data: body,
    ...(options || {}),
  })
}

/** POST /history/admin/delete */
export async function deleteHistoryMessageByAdmin(
  body: API.DeleteRequest,
  options?: { [key: string]: any },
) {
  return request<API.BaseResponseBoolean>('/history/admin/delete', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    data: body,
    ...(options || {}),
  })
}

/** GET /history/${historyId}/messages */
export async function listHistoryMessages(
  params: API.listHistoryMessagesParams,
  options?: { [key: string]: any },
) {
  const { historyId: param0, ...queryParams } = params
  return request<API.BaseResponsePageHistoryMessage>(`/history/${param0}/messages`, {
    method: 'GET',
    params: {
      pageSize: '10',
      ...queryParams,
    },
    ...(options || {}),
  })
}
