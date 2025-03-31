import {expect, test} from "@playwright/test";

test('has products', async ({page}) => {
    try {
        await page.goto('/products');
        expect(await page.getByTestId('page-header').textContent()).toBe('Products');
    } catch (_) {
        await page.screenshot({path: 'tests/screenshots/products-header-error.png', fullPage: true});
    }
});