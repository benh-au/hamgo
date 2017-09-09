import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { By }              from '@angular/platform-browser';
import { DebugElement }    from '@angular/core';

import { TopNavComponent } from './top-nav.component';

describe('TopNavComponent', () => {
  let component: TopNavComponent;
  let fixture: ComponentFixture<TopNavComponent>;
  let nav: DebugElement;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TopNavComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TopNavComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display a navigation bar', () => {
    fixture.detectChanges();
    nav = fixture.debugElement.query(By.css("nav"));
    expect(nav).not.toBeNull();
  });
});
