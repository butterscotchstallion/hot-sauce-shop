import {IBoard} from "./IBoard.ts";
import {IUser} from "../../User/types/IUser.ts";

export interface IBoardDetails {
    board: IBoard;
    moderators: IUser[];
    numBoardMembers: number;
    totalPosts: number;
}