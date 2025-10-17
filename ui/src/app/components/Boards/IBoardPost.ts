export interface IBoardPost {
    id: string;
    title: string;
    thumbnailFilename: string;
    postText: string;
    createdAt: Date;
    updatedAt: Date;
    boardId: number;
    boardName: string;
    boardSlug: string;
    parentId: number;
    createdAtUserId: number;
    createdByUsername: string;
    createdByUserSlug: string;
    slug: string;
}