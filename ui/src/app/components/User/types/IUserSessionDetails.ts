import {IUser} from "./IUser.ts";
import {IUserRole} from "./IUserRole.ts";

export interface IUserSessionDetails {
    user: IUser;
    roles: IUserRole[];
}