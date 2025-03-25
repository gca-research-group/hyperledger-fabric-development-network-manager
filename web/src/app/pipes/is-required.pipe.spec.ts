import { IsRequiredPipe } from './is-required.pipe';

describe('IsRequiredPipe', () => {
  it('create an instance', () => {
    const pipe = new IsRequiredPipe();
    expect(pipe).toBeTruthy();
  });
});
