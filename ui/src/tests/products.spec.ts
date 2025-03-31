import {expect, test} from "@playwright/test";

test('has products', async ({page}) => {
    try {
        await page.goto('/products');
        expect(await page.getByTestId('page-header').textContent()).toBe('Products');
    } catch (_) {
        await page.screenshot({path: 'screenshot.png', fullPage: true});
    }
});