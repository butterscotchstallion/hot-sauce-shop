import {z} from "zod";

export const ProductSchema = z.object({
    name: z.string().min(3).max(255),
    price: z.number().min(0.01).max(999999.99),
    description: z.string().min(3).max(1000000),
    shortDescription: z.string().min(3).max(1000),
    spiceRating: z.number().min(1).max(5),
})