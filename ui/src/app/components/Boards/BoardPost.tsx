import dayjs from "dayjs";
import {Card} from "primereact/card";
import {NavLink} from "react-router";
import TimeAgo from "react-timeago";
import {IBoardPost} from "./IBoardPost.ts";
import {Button} from "primereact/button";

interface IBoardPostProps {
    post: IBoardPost;
}

export default function BoardPost({post}: IBoardPostProps) {
    const createdAtFormatted = dayjs(post.createdAt).format('MMM DD, YYYY hh:mm A');
    let postImagePath: string = '/images/hot-pepper.png';
    if (post.thumbnailFilename) {
        postImagePath = `/images/posts/${post.thumbnailFilename}`;
    }
    return (
        <Card key={`post-${post.id}`}
              title={<NavLink to={`/boards/${post.boardSlug}/posts/${post.slug}`}>{post.title}</NavLink>}
              className="mb-4">
            <div>
                <div className="flex flex-column">
                    <img src={postImagePath} alt={post.title}/>
                    <div className="ml-6 min-h-[2rem]">
                        {post.postText}
                    </div>
                </div>
            </div>
            <div className="mt-6">
                <ul className="flex flex-wrap items-center text-gray-900 dark:text-white">
                    <li className="inline-block w-[120px]">
                        <div className="w-full flex justify-between items-center">
                            <Button icon={'pi pi-thumbs-up'} title="This is a high-quality post"/>
                            {post.voteSum}
                            <Button icon={'pi pi-thumbs-down'} title="This is a low-quality post"/>
                        </div>
                    </li>
                    <li className="ml-4 mr-4">&bull;</li>
                    <li>
                        <NavLink to={`/boards/${post.boardSlug}`}>
                            <i className="pi pi-list mr-1"></i> {post.boardName}
                        </NavLink>
                    </li>
                    <li className="ml-4 mr-4">&bull;</li>
                    <li>
                        <NavLink to={`/users/${post.createdByUserSlug}`}>
                            <i className="pi pi-user mr-1"></i> {post.createdByUsername}
                        </NavLink>
                    </li>
                    <li className="ml-4 mr-4">&bull;</li>
                    <li className="cursor-help">
                        <TimeAgo date={post.createdAt} title={createdAtFormatted}/>
                    </li>
                </ul>
            </div>
        </Card>
    )
}