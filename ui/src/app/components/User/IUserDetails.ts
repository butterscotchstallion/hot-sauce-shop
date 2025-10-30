import {IUser} from "./IUser.ts";
import {IUserRole} from "./IUserRole.ts";
import {IBoard} from "../Boards/IBoard.ts";
import {IUserLevelInfo} from "./IUserLevelInfo.ts";

export interface IUserDetails {
    user: IUser;
    roles: IUserRole[];
    userPostCount: number;
    postVoteSum: number;
    userPostVoteSum: number;
    userModeratedBoards: IBoard[];
    userLevelInfo: IUserLevelInfo;
}