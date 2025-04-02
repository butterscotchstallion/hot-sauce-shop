export interface IReview {
    id: string;
    title: string;
    comment: string;
    rating: number;
    createdAt: Date;
    updatedAt: Date;
    spiceRating: number;
    inventoryItemId: number;
    userId: number;
    username: string;
    userAvatarFilename: string;
    usernameSlug: string;
}