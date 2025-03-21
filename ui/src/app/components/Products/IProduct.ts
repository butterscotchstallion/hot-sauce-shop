export interface IProduct {
    id?: number;
    name: string;
    price: number;
    slug: string;
    description: string;
    shortDescription: string;
    createdAt?: Date;
    updatedAt?: Date;
    spiceRating: number;
}