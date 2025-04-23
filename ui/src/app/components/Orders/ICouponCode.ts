export interface ICouponCode {
    code: string;
    description: string;
    expirationDate: Date;
    createdAt: Date;
    updatedAt: Date;
    reductionPercent: number;
    couponTypeName: string;
}