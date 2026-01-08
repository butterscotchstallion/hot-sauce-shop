export interface IBoardDetailsPayload {
    isVisible: boolean;
    isPrivate: boolean;
    isPostApprovalRequired: boolean;
    isOfficial: boolean;
    minKarmaRequiredToPost: number;
    description: string;
    thumbnailFilename: string;
}