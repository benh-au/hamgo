import { HamgoPage } from './app.po';

describe('hamgo App', () => {
  let page: HamgoPage;

  beforeEach(() => {
    page = new HamgoPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!');
  });
});
