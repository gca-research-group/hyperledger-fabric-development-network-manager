import { ComponentFixture, TestBed } from '@angular/core/testing';

import { appConfig } from '@app/__tests__/app.config';

import { PeersSelectorComponent } from './peers-selector.component';

describe('PeersSelectorComponent', () => {
  let component: PeersSelectorComponent;
  let fixture: ComponentFixture<PeersSelectorComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [PeersSelectorComponent],
      providers: [appConfig.providers],
    }).compileComponents();

    fixture = TestBed.createComponent(PeersSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
