import {IUser} from "./IUser.ts";
import {IUserRole} from "./IUserRole.ts";

export interface IUserDetails {
    user: IUser;
    roles: IUserRole[];
    userPostCount: number;
    postVoteSum: number;
    userPostVoteSum: number;
}