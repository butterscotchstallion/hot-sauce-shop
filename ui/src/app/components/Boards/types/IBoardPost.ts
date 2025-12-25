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
    boardIsOfficial: boolean;
    parentId: number;
    createdAtUserId: number;
    createdByUsername: string;
    createdByUserSlug: string;
    slug: string;
    voteSum: number;
    isPinned: boolean;
    thumbnailWidth: number;
    thumbnailHeight: number;
}