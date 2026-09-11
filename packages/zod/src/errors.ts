import * as z from 'zod';

export const ZActionType = z.enum(['redirect']);

export const ZAction = z.object({
    type: ZActionType,
    message: z.string(),
    value: z.string(),
});

export const ZFieldError = z.object({
    field: z.string(),
    error: z.string(),
});

export const ZHTTPStatus = z.number().int().min(100).max(599);

export const ZAppError = z.object({
    action: ZAction,
    errors: ZFieldError,
    code: z.string(),
    message: z.string(),
    status: ZHTTPStatus,
});

export const ZErrorInfo = z.object({
    code: z.string(),
    message: z.string(),
    errors: z.array(ZFieldError),
    action: ZAction,
});

export const ZMeta = z.object({
    page: z.number(),
    limit: z.number(),
    total: z.number(),
    totalPages: z.number()
})

export const ZResponse = z.object({
    success: z.boolean(),
    status: ZHTTPStatus,
    data : z.any(),
    error: ZErrorInfo,
    meta: ZMeta
});

export type ActionType = z.infer<typeof ZActionType>
export type Action = z.infer<typeof ZAction>
export type FieldError = z.infer<typeof ZFieldError>
export type HTTPStatus = z.infer<typeof ZHTTPStatus>
export type ErrorInfo = z.infer<typeof ZErrorInfo>
export type Meta = z.infer<typeof ZMeta>
export type Response = z.infer<typeof ZResponse>