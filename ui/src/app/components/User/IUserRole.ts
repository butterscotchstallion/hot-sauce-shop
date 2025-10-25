import {TagProps} from "primereact/tag";

export interface IUserRole {
    id: number;
    name: string;
    slug: string;
    createdAt: string;
    updatedAt: string;
    colorClass: TagProps["severity"];
}