import { NetStatusPage } from './app.po';

describe('net-status App', () => {
  let page: NetStatusPage;

  beforeEach(() => {
    page = new NetStatusPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
