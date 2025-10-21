export interface IBoardPost {
    id: number;
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
    voteSum: number;
}