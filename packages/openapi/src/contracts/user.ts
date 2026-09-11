import type {ZodOpenApiPathsObject} from "zod-openapi";
import {ZResponse} from "@openmat/zod";


export const userPaths: ZodOpenApiPathsObject = {
    '/users': {
        get:{
            tags:["Users"],
            summary: "Fetch all users",
            description: "fetch a list of all users paginated response",
            responses: {
                '200': {
                    description: "users fetched successfully",
                    content: {
                        "application/json": {
                            schema: ZResponse
                        }
                    }
                },
                '400': {
                    description: "invalid request",
                    content: {
                        "application/json": {
                            schema: ZResponse
                        }
                    }
                },
                '401' : {
                    description: "unauthorized",
                    content: {
                        "application/json":{
                            schema: ZResponse
                        }
                    }
                }

            }
        }
    },
    '/users/{id}': {
        get: {
            tags: ["Users"],
            summary: "Fetch by ID",
            description: "fetch user by specified ID",
            responses: {
                '200': {
                    description: "user fetched successfully",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
                '400':  {
                    description: "invalid request",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
                '404':  {
                    description: "user not found",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
                '401':  {
                    description: "unauthorized",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
            }
        },
        patch: {
            tags: ["Users"],
            summary: "Update user",
            description: "update user by specified ID",
            responses: {
                '200': {
                    description: "user updated successfully",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
                '400':  {
                    description: "invalid request",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
                '404':  {
                    description: "user not found",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
                '401':  {
                    description: "unauthorized",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
            }
        },
        delete: {
            tags: ["Users"],
            summary: "Delete user",
            description: "delete user by specified ID",
            responses: {
                '200': {
                    description: "user deleted successfully",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
                '400':  {
                    description: "invalid request",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
                '404':  {
                    description: "user not found",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
                '401':  {
                    description: "unauthorized",
                    content: {
                        "application/json" : {
                            schema: ZResponse
                        }
                    }
                },
            }
        }
    }
}