export interface IBoard {
    id: number;
    displayName: string;
    slug: string;
    description: string;
    thumbnailFilename: string;
    createdAt: Date;
    updatedAt: Date;
    createdByUsername: string;
    createdAtByUserId: number;
    createdByUserSlug: string;
    visible: boolean;
}
