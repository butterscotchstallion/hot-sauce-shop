import dayjs from "dayjs";
import {Card} from "primereact/card";
import {NavLink} from "react-router";
import TimeAgo from "react-timeago";
import {IBoardPost} from "./IBoardPost.ts";
import {Button} from "primereact/button";
import {useEffect, useRef, useState} from "react";
import {Subject} from "rxjs";
import {addUpdateVote} from "./VoteService.ts";
import {VoteValue} from "./VoteValue.ts";

interface IBoardPostProps {
    post: IBoardPost;
}

export default function BoardPost({post}: IBoardPostProps) {
    const createdAtFormatted = dayjs(post.createdAt).format('MMM DD, YYYY hh:mm A');
    const [hasUpVoted, setHasUpVoted] = useState<boolean>(false);
    const [hasDownVoted, setHasDownVoted] = useState<boolean>(false);
    const vote$ = useRef<Subject<number>>(null);

    let postImagePath: string = '/images/hot-pepper.png';
    if (post.thumbnailFilename) {
        postImagePath = `/images/posts/${post.thumbnailFilename}`;
    }

    const onUpvoteClicked = () => {
        setHasUpVoted(true);
        vote(VoteValue.Upvote);
    }

    const onDownvoteClicked = () => {
        setHasDownVoted(true);
        vote(VoteValue.Downvote);
    }

    const vote = (voteValue: number) => {
        vote$.current = addUpdateVote(post.id, voteValue);
        vote$.current.subscribe({
            next: (_) => {
                setHasUpVoted(VoteValue.Upvote === voteValue);
                setHasDownVoted(VoteValue.Downvote === voteValue);
            },
            error: (err: string) => {
                setHasUpVoted(false);
                setHasDownVoted(false);
                console.log('Error voting: ' + err);
            }
        });
    }

    useEffect(() => {
        return () => {
            if (vote$.current) {
                vote$.current.unsubscribe();
            }
        }
    }, []);

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
                            <Button
                                icon={hasUpVoted ? 'pi pi-thumbs-up-fill' : 'pi pi-thumbs-up'}
                                title="This is a high-quality post"
                                disabled={hasUpVoted}
                                onClick={() => onUpvoteClicked()}
                            />
                            {post.voteSum}
                            <Button
                                icon={hasDownVoted ? 'pi pi-thumbs-down-fill' : 'pi pi-thumbs-down'}
                                title="This is a low-quality post"
                                disabled={hasDownVoted}
                                onClick={() => onDownvoteClicked()}
                            />
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