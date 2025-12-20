import {IBoardPost} from "./IBoardPost.ts";

export interface IBoardPostsResponse {
    posts: IBoardPost[];
    totalPosts: number;
}