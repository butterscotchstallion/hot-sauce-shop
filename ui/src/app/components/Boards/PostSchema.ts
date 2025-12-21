import {z} from "zod";

export const PostSchema = z.object({
    title: z.string().min(10, {message: "Title must be between 10 and 150 characters"}).max(150),
    postText: z.string().min(10, {message: "Post text must be between 10 and 150 characters"}).max(10000),
});