import { ComponentFixture, TestBed } from '@angular/core/testing';

import { PeersSelectorComponent } from './peers-selector.component';

describe('PeersSelectorComponent', () => {
  let component: PeersSelectorComponent;
  let fixture: ComponentFixture<PeersSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PeersSelectorComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(PeersSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
