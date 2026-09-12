import {
    ZCreateUserRequest,
    ZDeleteUserRequest,
    ZGetUserIdParams,
    ZResponse,
    ZUpdateUserRequest,
} from '@openmat/zod';
import type { ZodOpenApiPathsObject } from 'zod-openapi';

export const userPaths: ZodOpenApiPathsObject = {
    '/users': {
        post: {
            tags: ['Users'],
            summary: 'Create user',
            description: 'create a platform user',
            requestBody: {
                content: {
                    'application/json': {
                        schema: ZCreateUserRequest,
                    },
                },
            },
            responses: {
                201: {
                    description: 'user created successfully',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                400: {
                    description: 'invalid request',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                409: {
                    description: 'user already exists',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                404: {
                    description: 'user not found',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
            },
        },
        get: {
            tags: ['Users'],
            summary: 'Fetch all users',
            description: 'fetch a list of all users paginated response',
            responses: {
                '200': {
                    description: 'users fetched successfully',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '400': {
                    description: 'invalid request',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '401': {
                    description: 'unauthorized',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
            },
        },
    },
    '/users/{id}': {
        get: {
            tags: ['Users'],
            summary: 'Fetch by ID',
            description: 'fetch user by specified ID',
            requestParams: {
                path: ZGetUserIdParams,
            },
            responses: {
                '200': {
                    description: 'user fetched successfully',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '400': {
                    description: 'invalid request',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '404': {
                    description: 'user not found',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '401': {
                    description: 'unauthorized',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
            },
        },
        patch: {
            tags: ['Users'],
            summary: 'Update user',
            description: 'update user by specified ID',
            requestBody: {
                content: {
                    'application/json': {
                        schema: ZUpdateUserRequest,
                    },
                },
            },
            responses: {
                '200': {
                    description: 'user updated successfully',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '400': {
                    description: 'invalid request',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '404': {
                    description: 'user not found',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '401': {
                    description: 'unauthorized',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
            },
        },
        delete: {
            tags: ['Users'],
            summary: 'Delete user',
            description: 'delete user by specified ID',
            requestBody: {
                content: {
                    'application/json': {
                        schema: ZDeleteUserRequest,
                    },
                },
            },
            responses: {
                '200': {
                    description: 'user deleted successfully',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '400': {
                    description: 'invalid request',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '404': {
                    description: 'user not found',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
                '401': {
                    description: 'unauthorized',
                    content: {
                        'application/json': {
                            schema: ZResponse,
                        },
                    },
                },
            },
        },
    },
};
