import * as z from 'zod';

export const ZRole = z.enum(['user', 'student', 'instructor', 'admin']);

export const ZUser = z.object({
    id: z.string(),
    email: z.email(),
    password: z.string(),
    role: ZRole,
    createdAt: z.iso.datetime(),
    updatedAt: z.iso.datetime(),
});

export const ZUserResponse = z.object({
    id: z.string(),
    email: z.email(),
    role: ZRole,
    createdAt: z.iso.datetime(),
    updatedAt: z.iso.datetime(),
});

export const ZCreateUserRequest = z
    .object({
        email: z.email(),
        password: z
            .string()
            .min(8)
            .regex(/[A-Z]/, 'password must contain an uppercase letter')
            .regex(/[0-9]/, 'password must contain a number'),
    })
    .meta({
        example: {
            email: 'john@example.com',
            password: 'Password@123',
        },
    });

export const ZUpdateUserRequest = z
    .object({
        email: z.email().optional(),
        password: z.string().optional(),
    })
    .meta({
        example: {
            email: 'john@example.com',
            password: 'Password@123',
        },
    });

export const ZDeleteUserRequest = z
    .object({
        id: z.string(),
        email: z.email(),
    })
    .meta({
        example: {
            id: 'zd3fw3sfa4F',
            email: 'Password@123',
        },
    });

export const ZGetUserIdParams = z
    .object({
        id: z.string(),
    })
    .meta({
        example: {
            id: 'zd3fw3sfa4F',
        },
    });
